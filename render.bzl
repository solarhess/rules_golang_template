def _render_impl(ctx):
    render = ctx.executable._render
    template = ctx.file.template
    values = ctx.attr.literal_values
    out = ctx.outputs.out
    json_data_out = ctx.outputs.json_data_out

    file_data_values = ctx.attr.file_data_values
    json_data_values = ctx.attr.json_data_values

    file_input_args = []
    file_input_files = []

    #TODO update file input files and args
    for label,key in file_data_values.items() : 
        print("file input: ", label)
        file_arg = ctx.expand_location("-file=%s:$(location %s)" % (key, label.file.path), targets=[label])
        file_input_args.append(file_arg)
        file_input_files.append(label)
        print(file_arg)

    for label,key in json_data_values.items() : 
        file_arg = ctx.expand_location("-json=%s:$(location %s)" % (key, label.file.path), targets=[label])
        file_input_args.append(file_arg)
        file_input_files.append(label)


    outputs = [out]

    values_json = ctx.actions.declare_file(ctx.label.name + '.values.json')
    outputs.append(values_json)

    ctx.actions.write(values_json, struct(**values).to_json())
 
    ctx.actions.run(
        mnemonic = "RenderTemplate",
        inputs = [
            values_json,
            template,
        ] + file_input_files,
        executable=render,
        tools=[render],
        arguments=[
            "-data="+values_json.path, 
            "-template="+template.path, 
            "-output="+out.path,
            "-output-data-json="+json_data_out.path]
            + file_input_args,
        outputs = [json_data_out, out],
    )
    return [DefaultInfo(files = depset(outputs))]


_render = rule(
    implementation = _render_impl,
    attrs = {
        "template": attr.label(
            allow_files = True,
            single_file = True,
            mandatory = True,
        ),
        "file_data_values": attr.label_keyed_string_dict(
            allow_files=True,
            mandatory=False,
            allow_empty=True,
            default = {}
        ),
        "json_data_values": attr.label_keyed_string_dict(
            allow_files=True,
            mandatory=False,
            allow_empty=True,
            default = {}
        ),
        "literal_values": attr.string_dict(
            allow_empty=True, 
            default={}, 
            doc='The values to apply to the template', 
            mandatory=True, 
            non_empty=False
        ),
        # The label to the crd definition 'hybrises.modelt.hybris.com'
        "_render": attr.label(
            default = Label("//render:render"),
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "out": "%{name}",
        "json_data_out":"%{name}.data.json"
    },
)

def render(**kwargs):
    if "file_data_values" in kwargs : 
        kwargs["file_data_values"] = dict([[v,k] for k,v in kwargs["file_data_values"].items()])

    if "json_data_values" in kwargs : 
        kwargs["json_data_values"] = dict([[v,k] for k,v in kwargs["json_data_values"].items()])
    _render(**kwargs)

