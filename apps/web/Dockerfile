FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat wget

# Install dependencies
FROM base AS deps
WORKDIR /app

# Copy ALL package files (let npm figure out the workspace structure)
COPY package*.json ./
COPY apps/ ./apps/
COPY packages/ ./packages/
COPY turbo.json* ./

# Install everything
RUN npm install

# Build stage
FROM base AS builder
WORKDIR /app

# Copy everything from deps
COPY --from=deps /app/ ./

# Build (let npm handle workspace dependencies)
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Build the web app
RUN npm run build --workspace=apps/web

# Production stage
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy Next.js output
COPY --from=builder /app/apps/web/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT=3000

HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]
