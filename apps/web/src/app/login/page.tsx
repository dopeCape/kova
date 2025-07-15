"use client"
import React, { useState } from 'react';
import { Eye, EyeOff, Zap, Code2 } from 'lucide-react';

const DeploymentManager = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleAuth = () => {
  };


  return (
    <div className="min-h-screen bg-zinc-950 relative overflow-hidden">
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-950/30 via-transparent to-emerald-950/15"></div>
        <div
          className="absolute inset-0 opacity-5"
          style={{
            backgroundImage: `radial-gradient(circle at 50% 50%, #10b981 1px, transparent 1px)`,
            backgroundSize: '100px 100px'
          }}
        ></div>
        <div className="absolute top-1/3 right-1/4 w-96 h-96 bg-emerald-500/3 rounded-full blur-3xl"></div>
      </div>
      <div className="relative z-10 min-h-screen flex">
        <div className="flex-1 flex items-center justify-center p-12">
          <div className="w-full max-w-sm">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-white mb-2">
                Welcome Back
              </h2>
              <div className="flex items-center justify-center space-x-2 text-sm text-zinc-400 font-mono">
                <Code2 className="w-4 h-4" />
                <span>authenticate</span>
              </div>
            </div>

            <div className="space-y-5">
              <div className="relative group">
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-4 py-4 bg-zinc-900/80 backdrop-blur-sm border border-zinc-800 rounded-xl text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 focus:bg-zinc-900 transition-all font-mono text-sm group"
                  placeholder="email@domain.com"
                />
                <div className="absolute inset-0 rounded-xl bg-emerald-500/5 opacity-0 focus-within:opacity-100 transition-opacity pointer-events-none"></div>
              </div>

              <div className="relative group">
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-4 py-4 bg-zinc-900/80 backdrop-blur-sm border border-zinc-800 rounded-xl text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 focus:bg-zinc-900 transition-all font-mono text-sm pr-12"
                  placeholder="••••••••••••"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 top-1/2 transform -translate-y-1/2 text-zinc-500 hover:text-emerald-500 transition-colors"
                >
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
                <div className="absolute inset-0 rounded-xl bg-emerald-500/5 opacity-0 focus-within:opacity-100 transition-opacity pointer-events-none"></div>
              </div>

              <button
                onClick={handleAuth}
                className="w-full bg-emerald-500 hover:bg-emerald-400 text-white font-bold py-4 px-4 rounded-xl transition-all transform hover:scale-105 text-sm relative overflow-hidden group shadow-lg shadow-emerald-500/25"
              >
                <div className="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
                <span className="relative flex items-center justify-center space-x-2">
                  <Zap className="w-4 h-4" />
                  <span>Login</span>
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeploymentManager;
