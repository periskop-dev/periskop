## Build web
FROM node:lts AS fe-builder

ENV PORT 8080
ENV SERVER_URL localhost

WORKDIR /periskop-modules
COPY package-lock.json .
COPY package.json . 
RUN npm install
