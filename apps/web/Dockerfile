# Monorepo-aware Dockerfile for Next.js dashboard
FROM node:22-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy package files for the entire monorepo (root level)
COPY package*.json ./
COPY turbo.json* ./

# Copy workspace package.json files
COPY apps/web/package*.json ./apps/web/
COPY packages/ui/package*.json ./packages/ui/
COPY packages/*/package*.json ./packages/

# Install all dependencies using npm workspaces
RUN npm ci

# Builder stage
FROM base AS builder
WORKDIR /app

# Copy all dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages ./packages
COPY --from=deps /app/apps/web/node_modules ./apps/web/node_modules

# Copy source code for packages (UI components, etc.)
COPY packages/ ./packages/

# Copy source code for web app
COPY apps/web/ ./apps/web/

# Copy root level config files
COPY package*.json ./
COPY turbo.json* ./

# Build packages first (like packages/ui)
RUN npm run build --workspace=packages/ui --if-present

# Set environment variables for build
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Build the web application
RUN npm run build --workspace=apps/web

# Production stage
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy public assets
COPY --from=builder /app/apps/web/public ./public

# Copy Next.js standalone build
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]

