# File Template Rule for Bazel

This rule for Bazel allows you to render a template file from a Bazel build rule.
To use:

In your WORKSPACE add the following to import this repository rule to import this golang_rules_template project

```python
git_repository(
    name = "rules_golang_template",
    commit = "<master-branch-head-id>", 
    remote = "https://github.com/solarhess/rules_golang_template.git",
)
```


In your BUILD.bzl file, then you can then declare rules that run the template on an input file;

```python
# Load the rule at the top of the build file
load ("@rules_golang_template//:rules.bzl", "golang_template")
```

Then declare a template rule. You must declare the name and extension.
```python
golang_template(
    name="file-data-output",
    extension="txt",
    template="sample.tmpl.txt",

    # Explicitly set map values from bazel context variables
    literal_values={
        "sweaters": "15",
        "build_variable" : bazel_build_variable,
        "list" : "one",
    },
)
```

Finally use the template in the output\
```python
# Use a template output as input to another rule
sh_test(
    name = "test-file-data-output",
    size = "small",
    srcs = ["test-file-data-output.sh"],
    args = ["$(location :file-data-output.txt)"],
    data = [":file-data-output.txt"]
)
```

## Loading template context values from files

In addition to setting literal values from your build, you can also load values from files.
This can be very useful if you need to grab a value from a file produced by one rule and
render it into a template. For example, after getting credentials for an Azure SQL database,
you want to render a template containing the database password into a new YAML file.

The following code snippet in the [examples/BUILD.bazel] shows how to write a rule that loads
the text from file, and the object contents of json into the template rendering context.

```python
golang_template(
    name="file-data-output",
    extension="txt",
    template="sample.tmpl.txt",
    
    # Explicitly set map values from bazel context variables
    literal_values={
        "sweaters": "15",
        "build_variable" : bazel_build_variable,
        "list" : "one",
    },
    
    # Set the value of "script" to the contents of the file "script.txt"
    file_data_values = {
        "script" : ":script.txt"
    },

    # Set the value of "jsondata" to the json object parsed from the contents of jsondata.json 
    json_data_values = {
        "jsondata" : ":jsondata.json"
    }
)
```