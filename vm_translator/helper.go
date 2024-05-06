package vmtrans

import "text/template"

func loadTemplates(dir string) map[string]*template.Template {
	templates := make(map[string]*template.Template)
	for name, fn := range templateFileNames {
		t, err := template.New(fn).ParseFS(templateFiles, dir+"/"+fn)
		if err != nil {
			panic(err)
		}
		templates[name] = t
	}
	return templates
}
