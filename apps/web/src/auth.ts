import NextAuth from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"

declare module "next-auth" {
  interface Session {
    user: {
      id: string
      username: string
      email: string
    }
    accessToken: string
    refreshToken: string
  }


  interface User {
    id: string
    username: string
    email: string
    accessToken: string
    refreshToken: string
  }
}

declare module "@auth/core/jwt" {
  interface JWT {
    id: string
    username: string
    email: string
    accessToken: string
    refreshToken: string
  }

}



export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    CredentialsProvider({
      name: "credentials",
      credentials: {
        login: { label: "Email or Username", type: "text" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        if (!credentials?.login || !credentials?.password) {
          return null // Return null for missing credentials
        }

        try {
          const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
          const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              login: credentials.login,
              password: credentials.password,
            }),
          })

          if (!response.ok) {
            console.log('Login failed with status:', response.status);

            // For 401 (invalid credentials), return null instead of throwing
            if (response.status === 401) {
              return null
            }

            // For other errors, we can still throw
            const errorData = await response.json().catch(() => ({}))
            console.log('Error data:', errorData);

            // Return null for credential-related errors too
            if (errorData.code === 'INVALID_CREDENTIALS') {
              return null
            }

            // Only throw for server/network errors
            throw new Error(errorData.error || `Server error: ${response.status}`)
          }

          const data = await response.json()

          if (!data.user || !data.tokens) {
            console.log('Invalid response format:', data);
            return null
          }

          console.log('Login successful for user:', data.user.email);

          return {
            id: data.user.id,
            username: data.user.username,
            email: data.user.email,
            accessToken: data.tokens.access_token,
            refreshToken: data.tokens.refresh_token,
          }
        } catch (error) {
          console.error('Auth error:', error)

          // Handle network errors - throw these as they're not credential issues
          if (error instanceof TypeError && error.message.includes('fetch')) {
            throw new Error("Cannot connect to authentication server")
          }

          // For other errors, return null to avoid Configuration error
          return null
        }
      }
    })
  ],
  pages: {
    signIn: "/login",
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id
        token.username = user.username
        token.email = user.email
        token.accessToken = user.accessToken
        token.refreshToken = user.refreshToken
      }
      return token
    },
    async session({ session, token }) {
      session.user.id = token.id
      session.user.username = token.username
      session.user.email = token.email
      session.accessToken = token.accessToken
      session.refreshToken = token.refreshToken
      return session
    },
    async authorized({ auth, request: { nextUrl } }) {
      const isLoggedIn = !!auth?.user
      const isOnLogin = nextUrl.pathname === '/login'

      if (isOnLogin) {
        if (isLoggedIn) {
          return Response.redirect(new URL('/', nextUrl))
        }
        return true
      }

      return isLoggedIn
    },
  },
  session: {
    strategy: "jwt",
  },
  secret: process.env.AUTH_SECRET,
})
