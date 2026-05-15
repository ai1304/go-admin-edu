FROM node:22-alpine AS build
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-workspace.yaml ./
COPY packages/shared packages/shared
COPY apps/portal/package.json apps/portal/package.json
RUN pnpm --dir apps/portal install --no-frozen-lockfile
COPY apps/portal apps/portal
WORKDIR /app/apps/portal
RUN pnpm build

FROM nginx:1.27-alpine
COPY --from=build /app/apps/portal/dist /usr/share/nginx/html
COPY apps/portal/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
