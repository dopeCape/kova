import { NextResponse } from 'next/server'

export async function GET() {
  try {
    const health = {
      status: 'ok',
      service: 'kova-dashboard',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      environment: process.env.NODE_ENV || 'development',
      version: process.env.npm_package_version || '1.0.0',
    }


    return NextResponse.json(health, { status: 200 })
  } catch (error) {
    console.error('Health check failed:', error)

    return NextResponse.json(
      {
        status: 'error',
        service: 'kova-dashboard',
        timestamp: new Date().toISOString(),
        error: 'Health check failed',
      },
      { status: 503 }
    )
  }
}

export async function HEAD() {
  try {
    return new Response(null, { status: 200 })
  } catch (error) {
    return new Response(null, { status: 503 })
  }
}
