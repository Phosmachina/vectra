FROM node:latest

WORKDIR /vectra/static/css

RUN npm install -g postcss postcss-cli autoprefixer

CMD [ \
    "npx", \
    "postcss", "raw_style.css", \
    "--use", "autoprefixer", \
    "-o", "autoprefix_style.css" \
]
