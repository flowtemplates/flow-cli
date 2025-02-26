{# access value #}

{{ name }}

{# <---it's a comment hash syntax btw---> #}

{# available types #}

{{ false }} {# boolean #}

{{ 1 }} {# int #}

{{ -1.2 }} {# float #}

{{ "foo" }} {# string #}

{{ [1, 2] }} {# array #}

{{ {key: 1, another_key: 2} }} {# object #}


{# access object field #}

{{ user.name }}

{# access array index #}

{{ items[1] }}


{# operations evaluating #}

{{ my_num + 2 }}

{# bool results are not displayed as in jsx, which allows this to work #}

{{ flag && name }}


{# if/else statement #}

{% if some_num > 10 || some_num %}

{% else if some_num <= 2 && flag %}

{% else %}

{% endif %}


{# switch statement #}

{% switch some_var %}
  {% case "a" %}

  {% case "b" %}

  {% default %}

{% endif %}


{# ternary operator #}

{{ some_flag ? "this" : "or that" }}


{# filters #}

{{ "hello world" | snake }}  -> hello_world
{{ "hello world" | camel }}  -> helloWorld
{{ "hello world" | pascal }} -> HelloWorld
{{ "hello world" | kebab }}  -> hello-world
{{ "hello world" | title }}  -> Hello World
{{ "hello world" | upper }}  -> HELLO WORLD
{{ "hello world" | lower }}  -> hello world


{# loops #}

{% for item in items %}
  {{ item }}
{% endfor %}


{# with fancy range expression #}
{% for i in (1..4) %}
  {{ i }}
{% endfor %}

{# same as #}
{% for i in [1, 2, 3, 4] %}
  {{ i }}
{% endfor %}


{# variable declaration #}
{% let num = 4 %}
