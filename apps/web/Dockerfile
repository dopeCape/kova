FROM node:20-alpine AS base
RUN apk add --no-cache libc6-compat

# Dependencies stage
FROM base AS deps
WORKDIR /app

# Copy package files for the entire monorepo
COPY package*.json ./

# Copy workspace package.json files
COPY apps/web/package*.json ./apps/web/
COPY packages/ui/package*.json ./packages/ui/
COPY packages/typescript-config/package*.json ./packages/typescript-config/ || true
COPY packages/eslint-config/package*.json ./packages/eslint-config/ || true

# Copy any other config files
COPY turbo.json* ./

# Install all dependencies (npm workspaces hoists to root node_modules)
RUN npm ci

# Builder stage
FROM base AS builder
WORKDIR /app

# Copy ALL node_modules from deps (hoisted at root level)
COPY --from=deps /app/node_modules ./node_modules

# Copy package files again
COPY package*.json ./
COPY turbo.json* ./

# Copy source code for packages (shared components like UI)
COPY packages/ ./packages/

# Copy source code for web app
COPY apps/web/ ./apps/web/

# Build shared packages first (if they have build scripts)
RUN npm run build --workspace=packages/ui --if-present || echo "No UI package build script found"

# Set environment variables for Next.js build
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

# Copy Next.js standalone build output
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
