package excel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/config"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var (
	ErrInvalidFile      = errors.New("无效的Excel文件")
	ErrInvalidSheetName = errors.New("无效的工作表名称")
	ErrInvalidData      = errors.New("无效的数据")
)

// 配置项Excel导入导出处理

// ConfigItemExporter 配置项导出器
type ConfigItemExporter struct {
	ExportPath string
}

// NewConfigItemExporter 创建配置项导出器
func NewConfigItemExporter() *ConfigItemExporter {
	return &ConfigItemExporter{
		ExportPath: config.GetConfig().Excel.ExportPath,
	}
}

// Export 导出配置项到Excel
func (e *ConfigItemExporter) Export(ctx context.Context, items []models.ConfigurationItem) (string, error) {
	if len(items) == 0 {
		return "", ErrInvalidData
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("Failed to close Excel file", err)
		}
	}()

	// 设置工作表名
	sheetName := "配置项列表"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{
		"ID", "云服务商", "云产品", "配置项名称", "推荐配置值", 
		"风险说明", "检查方法", "配置方式", "参考资料",
		"创建时间", "更新时间",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置单元格样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Family: "微软雅黑",
			Size:   11,
			Color:  "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		logger.Error("Failed to create header style", err)
		return "", err
	}

	// 设置表头样式
	headerRange := fmt.Sprintf("A1:K1")
	if err := f.SetCellStyle(sheetName, "A1", "K1", headerStyle); err != nil {
		logger.Error("Failed to set header style", err)
		return "", err
	}

	// 填充数据
	for i, item := range items {
		rowIndex := i + 2 // 从第2行开始（第1行是表头）
		rowData := []interface{}{
			item.ID,
			item.Provider.Name,
			item.Product.Name,
			item.Name,
			item.RecommendedValue,
			item.RiskDescription,
			item.CheckMethod,
			item.ConfigurationMethod,
			item.Reference,
			item.CreatedAt.Format("2006-01-02 15:04:05"),
			item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		for j, cellData := range rowData {
			cell := fmt.Sprintf("%c%d", 'A'+j, rowIndex)
			f.SetCellValue(sheetName, cell, cellData)
		}
	}

	// 设置列宽
	colWidths := []float64{8, 15, 20, 40, 40, 30, 30, 30, 30, 20, 20}
	for i, width := range colWidths {
		col := string('A' + i)
		f.SetColWidth(sheetName, col, col, width)
	}

	// 确保导出目录存在
	if err := os.MkdirAll(e.ExportPath, 0755); err != nil {
		logger.Error("Failed to create export directory", err)
		return "", err
	}

	// 生成文件名
	filename := fmt.Sprintf("配置项列表_%s.xlsx", time.Now().Format("20060102150405"))
	filepath := filepath.Join(e.ExportPath, filename)

	// 保存文件
	if err := f.SaveAs(filepath); err != nil {
		logger.Error("Failed to save Excel file", err)
		return "", err
	}

	logger.Info("Excel file exported successfully", zap.String("filepath", filepath))
	return filepath, nil
}

// ConfigItemImporter 配置项导入器
type ConfigItemImporter struct {
	ImportPath string
}

// NewConfigItemImporter 创建配置项导入器
func NewConfigItemImporter() *ConfigItemImporter {
	return &ConfigItemImporter{
		ImportPath: config.GetConfig().Excel.ImportPath,
	}
}

// SaveUploadedFile 保存上传的Excel文件
func (i *ConfigItemImporter) SaveUploadedFile(file io.Reader, filename string) (string, error) {
	// 确保导入目录存在
	if err := os.MkdirAll(i.ImportPath, 0755); err != nil {
		logger.Error("Failed to create import directory", err)
		return "", err
	}

	// 生成文件路径
	filePath := filepath.Join(i.ImportPath, filename)

	// 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		logger.Error("Failed to create file", err, zap.String("filepath", filePath))
		return "", err
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error("Failed to copy file content", err)
		return "", err
	}

	logger.Info("File uploaded successfully", zap.String("filepath", filePath))
	return filePath, nil
}

// ImportConfigItems 从Excel文件导入配置项
func (i *ConfigItemImporter) ImportConfigItems(ctx context.Context, filePath string) ([]models.ConfigurationItem, error) {
	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		logger.Error("Failed to open Excel file", err, zap.String("filepath", filePath))
		return nil, ErrInvalidFile
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("Failed to close Excel file", err)
		}
	}()

	// 获取所有工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrInvalidSheetName
	}

	// 使用第一个工作表
	sheetName := sheets[0]

	// 获取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		logger.Error("Failed to get rows from Excel", err)
		return nil, err
	}

	if len(rows) < 2 { // 至少需要表头和一行数据
		return nil, ErrInvalidData
	}

	// 解析表头（第一行）
	headers := rows[0]
	// 验证表头...（这里简化处理，实际应用中可以更严格地验证表头）

	// 解析数据
	var items []models.ConfigurationItem
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 9 { // 至少需要9列数据
			logger.Warn("Skipping row due to insufficient columns", zap.Int("rowIndex", i+1))
			continue
		}

		// 提取数据
		// 注意：这里假设Excel中存储的是ID而不是名称，实际使用时可能需要转换
		providerID, err := strconv.ParseUint(row[1], 10, 32)
		if err != nil {
			logger.Warn("Invalid provider ID", zap.String("value", row[1]), zap.Int("rowIndex", i+1))
			continue
		}

		productID, err := strconv.ParseUint(row[2], 10, 32)
		if err != nil {
			logger.Warn("Invalid product ID", zap.String("value", row[2]), zap.Int("rowIndex", i+1))
			continue
		}

		item := models.ConfigurationItem{
			CloudProviderID:    uint(providerID),
			ProductID:          uint(productID),
			Name:               row[3],
			RecommendedValue:   row[4],
			RiskDescription:    row[5],
			CheckMethod:        row[6],
			ConfigurationMethod: row[7],
			Reference:          row[8],
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrInvalidData
	}

	return items, nil
}