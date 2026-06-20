import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'

interface AuthContextValue {
  username: string | null
  loading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

async function authRequest(path: string, options?: RequestInit) {
  const res = await fetch(`/api/v1/auth${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    let msg = '请求失败'
    try {
      const data = await res.json()
      if (data.error) msg = data.error
    } catch {}
    throw new Error(msg)
  }
  return res.json()
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [username, setUsername] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  const checkAuth = useCallback(async () => {
    try {
      const data = await authRequest('/me')
      setUsername(data.username)
    } catch {
      setUsername(null)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  const login = async (user: string, password: string) => {
    const data = await authRequest('/login', {
      method: 'POST',
      body: JSON.stringify({ username: user, password }),
    })
    setUsername(data.username)
  }

  const logout = async () => {
    try {
      await authRequest('/logout', { method: 'POST' })
    } finally {
      setUsername(null)
    }
  }

  return (
    <AuthContext.Provider value={{ username, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { username, loading } = useAuth()
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background text-on-surface-variant">
        <span className="material-symbols-outlined animate-spin text-[32px]">progress_activity</span>
      </div>
    )
  }
  if (!username) {
    window.location.href = '/login'
    return null
  }
  return <>{children}</>
}
