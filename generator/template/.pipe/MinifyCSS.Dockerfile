FROM tdewolff/minify:latest

WORKDIR /vectra/static/css

CMD [ \
    "minify", "--type", "css", \
    "autoprefix_style.css", "-o", "prod_style.css" \
]
