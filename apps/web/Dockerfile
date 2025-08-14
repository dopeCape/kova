FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat wget
WORKDIR /app

# Install dependencies
FROM base AS deps
COPY package*.json ./
COPY turbo.json* ./
COPY apps/web/package*.json ./apps/web/
COPY packages/ ./packages/
RUN npm install

# Build stage
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Build the web app
RUN npm run build --workspace=apps/web

# Production stage
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000

# Create user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy the standalone build (which includes the full apps/web structure)
COPY --from=builder /app/apps/web/.next/standalone ./

# Copy static files to match the structure
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/static ./apps/web/.next/static

# Set ownership
RUN chown -R nextjs:nodejs /app

USER nextjs
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# ~
# Use the working command from debug
CMD ["node", "apps/web/server.js"]
