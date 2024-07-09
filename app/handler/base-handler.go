package handler

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/config"
	"time"

	"github.com/labstack/echo/v4"
)

func NewBaseHandler(cfg config.Config) *BaseHandler {
	return &BaseHandler{cfg}
}

type BaseHandler struct {
	cfg config.Config
}

type Layouts struct {
	AppName      string
	Copyright    string
	LogoutPath   string
	LogoutMethod string
	SidebarMenus []*SidebarMenu
}

type SidebarMenu struct {
	GroupName string
	Items     []*SidebarMenuItem
}

type SidebarMenuItem struct {
	Key      string
	Html     string
	Href     string
	IsActive bool
	SubMenus []*SidebarMenuItem
}

func (h *BaseHandler) GetAuthor(e echo.Context) (*entity.User, error) {
	// get author
	_author := e.Get(h.cfg.JwtUserContextKey)
	author, ok := _author.(*entity.User)
	if !ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.BaseHandlerGetAuthorError, "failed type assertion author"),
		}
	}
	return author, nil
}

func (h *BaseHandler) GetLayouts(activeMenu string) *Layouts {
	return &Layouts{
		AppName:      config.AppName,
		Copyright:    h.GetCopyright(),
		LogoutPath:   entity.WebLogoutPath,
		LogoutMethod: entity.WebLogoutMethod,
		SidebarMenus: h.GetSidebarMenu(activeMenu),
	}
}

func (h *BaseHandler) GetCopyright() string {
	return fmt.Sprintf(
		`Copyright &copy; <a href="https://github.com/gendutski/rap-c" target="_blank">%s</a> %s`,
		config.AppName,
		time.Now().Format("2006"))
}

func (h *BaseHandler) GetSidebarMenu(activeMenu string) []*SidebarMenu {
	menus := []*SidebarMenu{
		// dashboard
		{
			GroupName: "",
			Items: []*SidebarMenuItem{{
				Key:  "dashboard",
				Html: `<i class="fas fa-fw fa-tachometer-alt"></i> <span>Dashboard</span></a>`,
				Href: "/dashboard",
			}},
		},
		// unit, ingredients, recipe, production and sale
		{
			GroupName: "Menu Utama",
			Items: []*SidebarMenuItem{
				{
					Key:  "unit",
					Html: `<i class="fa-solid fa-scale-balanced"></i> <span>Satuan Ukuran</span></a>`,
					Href: "/unit",
				},
				{
					Key:  "ingredients",
					Html: `<i class="fa-solid fa-fire"></i> <span>Bahan Baku</span>`,
					SubMenus: []*SidebarMenuItem{
						{
							Key:  "ingredients-list",
							Html: `<i class="fa-solid fa-clipboard-list"></i> <span>Daftar Bahan Baku</span>`,
							Href: "/ingredients",
						},
						{
							Key:  "update-stock",
							Html: `<i class="fa-solid fa-list-check"></i> <span>Pergerakan Stok</span>`,
						},
					},
				},
				{
					Key:  "recipe",
					Html: `<i class="fa-solid fa-cake-candles"></i> <span>Resep &amp; Produksi</span>`,
					SubMenus: []*SidebarMenuItem{
						{
							Key:  "recipe-list",
							Html: `<i class="fa-solid fa-clipboard-list"></i> <span>Daftar Resep</span>`,
							Href: "/recipe",
						},
						{
							Key:  "product",
							Html: `<i class="fa-solid fa-blender-phone"></i> <span>Daftar Produksi</span>`,
							Href: "/product",
						},
						{
							Key:  "sale",
							Html: `<i class="fa-solid fa-shop"></i> <span>Penjualan</span>`,
							Href: "/sale",
						},
					},
				},
			},
		},
		// Report
		{
			GroupName: "Laporan",
			Items: []*SidebarMenuItem{
				{
					Key:  "transaction",
					Html: `<i class="fa-solid fa-hand-holding-dollar"></i> <span>Transaksi</span>`,
					Href: "/transaction",
				},
			},
		},
		// Admin
		{
			GroupName: "Administrasi",
			Items: []*SidebarMenuItem{
				{
					Key:  "user",
					Html: `<i class="fa-regular fa-user"></i> <span>Daftar User</span>`,
					Href: "/user",
				},
				{
					Key:  "backup",
					Html: `<i class="fa-solid fa-cloud-arrow-up"></i> <span>Backup Data</span>`,
				},
			},
		},
	}

	for _, groupMenu := range menus {
		exists := false
		for _, menu := range groupMenu.Items {
			if menu.Key == activeMenu {
				exists = true
				menu.IsActive = true
				break
			}
			for _, subMenu := range menu.SubMenus {
				if subMenu.Key == activeMenu {
					exists = true
					menu.IsActive = true
					subMenu.IsActive = true
					break
				}
			}
			if exists {
				break
			}
		}
		if exists {
			break
		}
	}

	return menus
}
