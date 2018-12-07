// Package config provides a content type to manage the Ponzu system's configuration
// settings for things such as its name, domain, HTTP(s) port, email, server defaults
// and backups.
package config

import (
	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

// Config represents the confirgurable options of the system
type Config struct {
	item.Item

	Name                    string   `json:"name"`
	Domain                  string   `json:"domain"`
	BindAddress             string   `json:"bind_addr"`
	HTTPPort                string   `json:"http_port"`
	HTTPSPort               string   `json:"https_port"`
	AdminEmail              string   `json:"admin_email"`
	ClientSecret            string   `json:"client_secret"`
	Etag                    string   `json:"etag"`
	DisableCORS             bool     `json:"cors_disabled"`
	DisableGZIP             bool     `json:"gzip_disabled"`
	DisableHTTPCache        bool     `json:"cache_disabled"`
	CacheMaxAge             int64    `json:"cache_max_age"`
	CacheInvalidate         []string `json:"cache"`
	BackupBasicAuthUser     string   `json:"backup_basic_auth_user"`
	BackupBasicAuthPassword string   `json:"backup_basic_auth_password"`
}

const (
	dbBackupInfo = `
		<p class="flow-text">数据库备份凭证:</p>
		<p>添加用户名和密码，以通过HTTP下载你的数据备份</p>
	`
)

// String partially implements item.Identifiable and overrides Item's String()
func (c *Config) String() string { return c.Name }

// MarshalEditor writes a buffer of html to edit a Post and partially implements editor.Editable
func (c *Config) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(c,
		editor.Field{
			View: editor.Input("Name", c, map[string]string{
				"label":       "站点名称",
				"placeholder": "添加内部使用的站点名称",
			}),
		},
		editor.Field{
			View: editor.Input("Domain", c, map[string]string{
				"label":       "域名(用来请求SSL证书)",
				"placeholder": "例如www.example.com或example.com",
			}),
		},
		editor.Field{
			View: editor.Input("BindAddress", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("HTTPPort", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("HTTPSPort", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("AdminEmail", c, map[string]string{
				"label": "管理员邮箱(通知系统信息)",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"label":    "客户端密码(用于验证请求，不可分享)",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"label":    "Etag头部(用来缓存资源)",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableCORS", c, map[string]string{
				"label": "使CORS失效(只有" + c.Domain + "可以使用你的数据)",
			}, map[string]string{
				"true": "使CORS失效",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableGZIP", c, map[string]string{
				"label": "使GZIP失效(提高服务器速度和带宽)",
			}, map[string]string{
				"true": "使GZIP失效",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableHTTPCache", c, map[string]string{
				"label": "使HTTP Cache失效(重写'Cache-Control'头部)",
			}, map[string]string{
				"true": "使HTTP Cache失效",
			}),
		},
		editor.Field{
			View: editor.Input("CacheMaxAge", c, map[string]string{
				"label": "HTTP caching有效期(单位s, 0 = 2592000)",
				"type":  "text",
			}),
		},
		editor.Field{
			View: editor.Checkbox("CacheInvalidate", c, map[string]string{
				"label": "保存缓存失效",
			}, map[string]string{
				"invalidate": "缓存失效",
			}),
		},
		editor.Field{
			View: []byte(dbBackupInfo),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthUser", c, map[string]string{
				"label":       "基本HTTP认证用户",
				"placeholder": "输入基本认证用户名",
				"type":        "text",
			}),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthPassword", c, map[string]string{
				"label":       "基本HTTP认证密码",
				"placeholder": "输入基本认证密码",
				"type":        "password",
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	open := []byte(`
	<div class="card">
		<div class="card-content">
			<div class="card-title">系统配置</div>
		</div>
		<form action="/admin/configure" method="post">
	`)
	close := []byte(`</form></div>`)
	script := []byte(`
	<script>
		$(function() {
			// hide default fields & labels unnecessary for the config
			var fields = $('.default-fields');
			fields.css('position', 'relative');
			fields.find('input:not([type=submit])').remove();
			fields.find('label').remove();
			fields.find('button').css({
				position: 'absolute',
				top: '-10px',
				right: '0px'
			});

			var contentOnly = $('.content-only.__ponzu');
			contentOnly.hide();
			contentOnly.find('input, textarea, select').attr('name', '');

			// adjust layout of td so save button is in same location as usual
			fields.find('td').css('float', 'right');

			// stop some fixed config settings from being modified
			fields.find('input[name=client_secret]').attr('name', '');
		});
	</script>
	`)

	view = append(open, view...)
	view = append(view, close...)
	view = append(view, script...)

	return view, nil
}
