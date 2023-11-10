FROM node:latest

WORKDIR /vectra/static/css

RUN npm install -g sass

CMD [ \
    "sass", "--embed-source-map", \
    "style.sass", \
    "raw_style.css" \
]
