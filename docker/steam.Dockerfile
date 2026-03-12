FROM node:24.12.0-alpine3.22 AS build
WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY tsconfig.json tsconfig.build.json ./
COPY src ./src
RUN pnpm run build

RUN pnpm prune --prod

FROM node:24.12.0-alpine3.22 AS runtime
WORKDIR /app

ENV NODE_ENV=production
ENV APP_ENV=production

COPY --from=build /app/package.json ./package.json
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist

EXPOSE 3002

CMD ["node", "dist/index.js"]
