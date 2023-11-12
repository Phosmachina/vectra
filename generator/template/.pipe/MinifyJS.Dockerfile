FROM tdewolff/minify:latest

WORKDIR /vectra/static/js

CMD [ \
    "minify", "--type", "js", \
    "main.js", "-o", "prod_main.js" \
]
