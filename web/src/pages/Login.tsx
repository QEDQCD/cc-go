import { useState, type FormEvent } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../auth'

export default function Login() {
  const { username, loading, login } = useAuth()
  const [user, setUser] = useState('admin')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  if (!loading && username) {
    return <Navigate to="/" replace />
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      await login(user, password)
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen relative overflow-hidden bg-background text-on-surface">
      {/* Background */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_left,_rgba(78,222,163,0.12)_0%,_transparent_50%)]" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_right,_rgba(190,198,224,0.08)_0%,_transparent_50%)]" />
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage: 'linear-gradient(rgba(211,228,254,0.5) 1px, transparent 1px), linear-gradient(90deg, rgba(211,228,254,0.5) 1px, transparent 1px)',
          backgroundSize: '48px 48px',
        }}
      />

      <div className="relative z-10 min-h-screen flex items-center justify-center p-6">
        <div className="w-full max-w-4xl grid md:grid-cols-2 gap-8 items-center">
          {/* Intro */}
          <div className="space-y-6">
            <div className="flex items-center gap-3">
              <img src="/cc-go.png" alt="cc-go" className="w-12 h-12 rounded-xl" />
              <div>
                <h1 className="text-3xl font-bold text-on-surface tracking-tight">cc-go</h1>
                <p className="text-sm text-on-surface-variant">远程 Claude Code 控制台</p>
              </div>
            </div>

            <p className="text-on-surface-variant leading-relaxed text-[15px]">
              cc-go 是一款基于 Claude Code 的远程编码工具。通过微信机器人接管 Claude Code 会话，
              你可以在手机上批准权限请求、查看 AI 回复、启动/切换会话，真正做到随时随地编码。
            </p>

            <div className="grid grid-cols-2 gap-3">
              {[
                { icon: 'chat', label: '微信远程控制' },
                { icon: 'shield', label: '权限审批' },
                { icon: 'dashboard', label: 'Web 管理面板' },
                { icon: 'bolt', label: '实时推送' },
              ].map(item => (
                <div
                  key={item.label}
                  className="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-surface-container/60 border border-outline-variant/40"
                >
                  <span className="material-symbols-outlined text-secondary text-[20px]">{item.icon}</span>
                  <span className="text-sm text-on-surface-variant">{item.label}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Login form */}
          <div className="bg-surface-container/80 backdrop-blur-xl border border-outline-variant/50 rounded-2xl p-8 shadow-2xl">
            <h2 className="text-xl font-semibold text-on-surface mb-1">登录管理面板</h2>
            <p className="text-sm text-on-surface-variant mb-6">请输入账号密码以继续</p>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm text-on-surface-variant mb-1.5">用户名</label>
                <input
                  type="text"
                  value={user}
                  onChange={e => setUser(e.target.value)}
                  autoComplete="username"
                  className="w-full px-4 py-2.5 rounded-lg bg-surface-container-low border border-outline-variant text-on-surface placeholder:text-on-surface-variant/50 focus:outline-none focus:border-secondary/60 transition-colors"
                  placeholder="admin"
                />
              </div>
              <div>
                <label className="block text-sm text-on-surface-variant mb-1.5">密码</label>
                <input
                  type="password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  autoComplete="current-password"
                  className="w-full px-4 py-2.5 rounded-lg bg-surface-container-low border border-outline-variant text-on-surface placeholder:text-on-surface-variant/50 focus:outline-none focus:border-secondary/60 transition-colors"
                  placeholder="请输入密码"
                />
              </div>

              {error && (
                <div className="flex items-center gap-2 text-error text-sm px-3 py-2 rounded-lg bg-error/10 border border-error/20">
                  <span className="material-symbols-outlined text-[18px]">error</span>
                  {error}
                </div>
              )}

              <button
                type="submit"
                disabled={submitting || !user || !password}
                className="w-full py-2.5 rounded-lg bg-secondary text-on-secondary font-medium hover:bg-secondary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
              >
                {submitting ? (
                  <>
                    <span className="material-symbols-outlined animate-spin text-[20px]">progress_activity</span>
                    登录中...
                  </>
                ) : (
                  <>
                    <span className="material-symbols-outlined text-[20px]">login</span>
                    登录
                  </>
                )}
              </button>
            </form>

            <p className="mt-6 text-xs text-on-surface-variant/60 text-center">
              登录凭据见服务端 ~/.cc-go/auth.json（可参考 auth.example.json）
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
