package main

import (
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func main() {

	devices := DevicesInit()
	loadKeyNames()

	devicetypesfile := jen.NewFilePathName("../hadiscovery/devicetypes.go", "hadiscovery")

	devicetypesfile.ImportAlias("github.com/eclipse/paho.mqtt.golang", "mqtt")

	deviceinitfile := jen.NewFilePathName("../hadiscovery/deviceinit.go", "hadiscovery")

	devicefunctionsfile := jen.NewFilePathName("../hadiscovery/devicefunctions.go", "hadiscovery")

	sort.Strings(keyNames)

	for _, d := range devices {

		// d.GetRawID()
		devicefunctionsfile.Func().Params(
			jen.Id("d").Id(strcase.ToCamel(d.Name)),
		).Id("GetRawId").Params().String().Block(
			jen.Return(jen.Lit(d.Name)),
		)

		// d.GetUniqueID()
		if d.JSONContainer.Exists("unique_id") {
			devicefunctionsfile.Func().Params(
				jen.Id("d").Id(strcase.ToCamel(d.Name)),
			).Id("GetUniqueId").Params().String().Block(
				jen.Return(jen.Id("d.UniqueId")),
			)
		} else {
			devicefunctionsfile.Func().Params(
				jen.Id("d").Id(strcase.ToCamel(d.Name)),
			).Id("GetUniqueId").Params().String().Block(
				jen.Return(jen.Lit("")),
			)
		}

		st := make(map[string][]*jen.Statement)

		// Add standalone base level fields
		for _, key := range keyNames {
			if d.JSONContainer.Exists(key) {
				st[key] = append(st[key], d.FieldAdder(key), d.FunctionAdder(key))
			}
		}

		if d.JSONContainer.Exists("device") {
			st["device"] = append(st["device"], jen.Id(strcase.ToCamel("device")).Struct(
				jen.Id(strcase.ToCamel("configuration_url")).String().Tag(map[string]string{"json": "configuration_url"}),
				jen.Id(strcase.ToCamel("connections")).Index().String().Tag(map[string]string{"json": "connections"}),
				jen.Id(strcase.ToCamel("identifiers")).Index().String().Tag(map[string]string{"json": "identifiers"}),
				jen.Id(strcase.ToCamel("manufacturer")).String().Tag(map[string]string{"json": "manufacturer"}),
				jen.Id(strcase.ToCamel("model")).String().Tag(map[string]string{"json": "model"}),
				jen.Id(strcase.ToCamel("name")).String().Tag(map[string]string{"json": "name"}),
				jen.Id(strcase.ToCamel("suggested_area")).String().Tag(map[string]string{"json": "suggested_area"}),
				jen.Id(strcase.ToCamel("sw_version")).String().Tag(map[string]string{"json": "sw_version"}),
				jen.Id(strcase.ToCamel("via_device")).String().Tag(map[string]string{"json": "via_device"}),
			).Tag(map[string]string{"json": "device"}))
		}

		// d.PopulateDevice()
		if d.JSONContainer.Exists("device") {
			devicefunctionsfile.Func().Params(
				jen.Id("d").Id("*"+strcase.ToCamel(d.Name)),
			).Id("PopulateDevice").Params().Block(
				jen.Id("d.Device.Manufacturer").Op("=").Id("Manufacturer"),
				jen.Id("d.Device.Model").Op("=").Id("SoftwareName"),
				jen.Id("d.Device.Name").Op("=").Id("InstanceName"),
				jen.Id("d.Device.SwVersion").Op("=").Id("SWVersion"),
			)
		} else {
			devicefunctionsfile.Func().Params(
				jen.Id("d").Id("*" + strcase.ToCamel(d.Name)),
			).Id("PopulateDevice").Params().Block()
		}

		sortedKeys := []string{}
		for key := range st {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		// Device MQTT Struct
		devicetypesfile.Type().Id(strcase.ToCamel(d.Name)).StructFunc(
			func(g *jen.Group) {
				for _, key := range sortedKeys {
					v := st[key]
					for _, val := range v {
						g.Add(val)
					}
				}
				g.Id("RawId").String().Tag(map[string]string{"json": "-"})
			},
		)

		deviceinitfile.Func().Params(jen.Id("d").Id(strcase.ToCamel(d.Name))).Id("Init").Params().BlockFunc(func(g *jen.Group) {
			if d.JSONContainer.Exists("retain") {
				g.Add(jen.Id("d").Dot("Retain").Op("=").Lit(false))
			}
			g.Add(jen.Id("d").Dot("PopulateDevice").Params())
			g.Add(jen.Id("d").Dot("PopulateTopics").Params())
		})

		// d.PopulateTopics()
		deviceinitfile.Func().Params(
			jen.Id("d").Id("*" + strcase.ToCamel(d.Name)),
		).Id("PopulateTopics").Params().BlockFunc(func(g *jen.Group) {
			for _, name := range keyNames {
				if strings.HasSuffix(name, "_topic") && d.JSONContainer.Exists(name) {
					lName := strcase.ToCamel(strings.TrimSuffix(name, "_topic"))
					g.Add(
						jen.If(
							jen.Id("d").Dot(lName + "Func").Op("!=").Nil(),
						).BlockFunc(
							func(g *jen.Group) {
								g.Add(jen.Id("d").Dot(strcase.ToCamel(name)).Op("=").Id("GetTopic").Params(jen.Id("d"), jen.Lit(name)))
								if IsCommand(name, d) {
									g.Add(jen.Id("topicStore").Index(
										jen.Id("d").Dot(strcase.ToCamel(name)),
									).Op("=").Id("&d").Dot(lName + "Func"))
								}
							},
						),
					)
				}
			}
		})

	}

	devicetypesfile.Type().Id("Device").Interface(
		// jen.UnionFunc(
		// 	func(g *jen.Group) {
		// 		for _, d := range devices {
		// 			g.Add(jen.Id(strcase.ToCamel(d.Name)))
		// 		}
		// 	},
		// ),
		jen.Id("GetRawId").Params().String(),
		jen.Id("GetUniqueId").Params().String(),
		jen.Id("PopulateDevice").Params(),
		jen.Id("PopulateTopics").Params(),
	)

	devicetypesfile.Save("../hadiscovery/devicetypes.go")

	deviceinitfile.Save("../hadiscovery/deviceinit.go")

	devicefunctionsfile.Save("../hadiscovery/devicefunctions.go")

}
