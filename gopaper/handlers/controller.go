// base controller

package handlers

import (
	"appengine"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"gopaper/config"
	"gopaper/ex/templateex"
)

type Controller struct {
	Name        string   // name of the controller
	Action      string   // name of the action
	ViewName    string   // name of the view
	NestedViews []string // names of nested views

	Context appengine.Context // app engine context

	Response http.ResponseWriter
	Request  *http.Request
}

// by using reflection, create the requested controller and call it
func MakeHandler(action interface{}) http.HandlerFunc {
	return func(rsp http.ResponseWriter, req *http.Request) {

		// find the type of the controller and create a new one
		controllerType := reflect.TypeOf(action).In(0)
		newController := reflect.New(controllerType)

		// find the action of the controller
		actionName := runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
		actionName = actionName[strings.LastIndex(actionName, ".")+1 : len(actionName)]

		// fill the base controller
		baseController := new(Controller)
		baseController.Response = rsp
		baseController.Request = req
		baseController.Name = controllerType.Name()
		baseController.Context = appengine.NewContext(req)
		baseController.Action = actionName
		// inject the base controller
		newController.Elem().FieldByName("Controller").Set(reflect.ValueOf(baseController))

		// execute the action
		reflect.ValueOf(action).Call([]reflect.Value{newController.Elem()})
	}
}

// render the requested view
func (c *Controller) Render(data interface{}) {

	// if no view is requested then take the default
	if len(c.ViewName) == 0 {
		c.ViewName = c.Name + "/" + c.Action
	}

	// create a new template
	t := template.New("_base.html").Funcs(templateFuncs)

	// inject nested views
	for _, value := range c.NestedViews {
		t = template.Must(t.ParseFiles("gopaper/templates/" + value + ".html"))
	}

	// inject the base and requested view
	t = template.Must(t.ParseFiles(
		"gopaper/templates/_base.html",
		"gopaper/templates/"+c.ViewName+".html",
	))

	// wrap the requested data
	type Data interface{}
	wrapper := struct {
		Data
		ControllerName string
		ActionName     string
		Config         *config.Config
	}{
		data,
		c.Name,
		c.Action,
		config.Default,
	}

	// show it
	err := t.Execute(c.Response, wrapper)

	// if any error, show it only on dev server
	if err != nil && config.Default.IsDevApp {
		fmt.Fprintf(c.Response, err.Error())
	}
}

var (
	// template functions
	templateFuncs = template.FuncMap{
		"equal":    templateex.Equal,
		"plus":     templateex.Addition,
		"date":     templateex.FormatDate,
		"yesno":    templateex.FormatBool,
		"htmlSafe": templateex.HtmlSafe,
	}
)

// redirect to the requested path
func (c *Controller) Redirect(path string) {
	http.Redirect(c.Response, c.Request, path, http.StatusFound)
}
