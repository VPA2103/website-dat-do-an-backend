package ingest

import (
	"strconv"
	"strings"

	"github.com/vpa/quanlynhahang-backend/ai/core"
)

func menuItemToDocument(item core.MenuItem) string {
	parts := []string{"Món: " + item.Name}
	if item.Description != nil && strings.TrimSpace(*item.Description) != "" {
		parts = append(parts, "Mô tả: "+strings.TrimSpace(*item.Description))
	}
	if item.Price != nil {
		parts = append(parts, "Giá: "+strconv.FormatFloat(*item.Price, 'f', -1, 64))
	}
	if len(item.Tags) > 0 {
		parts = append(parts, "Tags: "+strings.Join(item.Tags, ", "))
	}
	if len(item.Allergens) > 0 {
		parts = append(parts, "Dị ứng: "+strings.Join(item.Allergens, ", "))
	}
	if len(item.Ingredients) > 0 {
		parts = append(parts, "Nguyên liệu: "+strings.Join(item.Ingredients, ", "))
	}
	return strings.Join(parts, ". ")
}

func restaurantToDocument(info core.RestaurantInfo) string {
	lines := []string{"Thông tin nhà hàng"}
	if info.Name != nil && *info.Name != "" {
		lines = append(lines, "name: "+*info.Name)
	}
	if info.Address != nil && *info.Address != "" {
		lines = append(lines, "address: "+*info.Address)
	}
	if info.OpenHours != nil && *info.OpenHours != "" {
		lines = append(lines, "open_hours: "+*info.OpenHours)
	}
	if info.Phone != nil && *info.Phone != "" {
		lines = append(lines, "phone: "+*info.Phone)
	}
	if info.Style != nil && *info.Style != "" {
		lines = append(lines, "style: "+*info.Style)
	}
	if info.Policies != nil && *info.Policies != "" {
		lines = append(lines, "policies: "+*info.Policies)
	}
	return strings.Join(lines, "\n")
}
