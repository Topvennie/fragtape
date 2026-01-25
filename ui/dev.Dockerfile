FROM node:24.12.0-alpine3.22

WORKDIR /app

RUN npm install -g pnpm@10.16.1

COPY package.json pnpm-lock.yaml ./

CMD [ "sh", "-c", "pnpm install && pnpm run host" ]
