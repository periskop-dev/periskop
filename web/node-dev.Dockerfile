## Build web
FROM node:lts AS fe-builder

ENV PORT 8080

WORKDIR /periskop-modules
COPY package-lock.json .
COPY package.json . 
RUN npm install
