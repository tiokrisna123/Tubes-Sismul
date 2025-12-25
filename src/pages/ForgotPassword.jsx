import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Heart, Mail, Lock, Loader2, ArrowLeft, CheckCircle, Eye, EyeOff } from 'lucide-react';
import axios from 'axios';
import './Auth.css';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const ForgotPassword = () => {
    const [email, setEmail] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        // Validate passwords match
        if (newPassword !== confirmPassword) {
            setError('Password baru dan konfirmasi password tidak sama');
            return;
        }

        // Validate password length
        if (newPassword.length < 6) {
            setError('Password minimal 6 karakter');
            return;
        }

        setLoading(true);

        try {
            await axios.post(`${API_URL}/auth/reset-password`, {
                email: email,
                new_password: newPassword
            });
            setSuccess(true);
        } catch (err) {
            setError(err.response?.data?.error || 'Gagal reset password. Pastikan email terdaftar.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-container">
            {/* Animated Background Elements */}
            <div className="auth-bg-elements">
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
            </div>

            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <Heart className="logo-icon" />
                        <span>Live for Health</span>
                    </div>
                    <h1>Reset Password</h1>
                    <p>Masukkan email dan password baru Anda</p>
                </div>

                {success ? (
                    <div className="auth-form">
                        <div className="success-message">
                            <CheckCircle size={32} style={{ marginBottom: 12 }} />
                            <p>Password berhasil direset!</p>
                            <p style={{ fontSize: 13, marginTop: 8, opacity: 0.8 }}>
                                Silakan login dengan password baru Anda.
                            </p>
                        </div>
                        <div className="back-to-login" style={{ marginTop: 24 }}>
                            <Link to="/login">
                                <ArrowLeft size={18} />
                                Login Sekarang
                            </Link>
                        </div>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="auth-form">
                        <div className="forgot-password-info">
                            <p>
                                🔐 Masukkan email yang terdaftar dan password baru Anda.
                            </p>
                        </div>

                        {error && <div className="auth-error">{error}</div>}

                        <div className="input-group">
                            <Mail className="input-icon" />
                            <input
                                type="email"
                                placeholder="Email terdaftar"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>

                        <div className="input-group">
                            <Lock className="input-icon" />
                            <input
                                type={showPassword ? "text" : "password"}
                                placeholder="Password baru"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                                minLength={6}
                            />
                            <button
                                type="button"
                                className="password-toggle"
                                onClick={() => setShowPassword(!showPassword)}
                                tabIndex={-1}
                            >
                                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                            </button>
                        </div>

                        <div className="input-group">
                            <Lock className="input-icon" />
                            <input
                                type={showConfirmPassword ? "text" : "password"}
                                placeholder="Konfirmasi password baru"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                                minLength={6}
                            />
                            <button
                                type="button"
                                className="password-toggle"
                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                tabIndex={-1}
                            >
                                {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                            </button>
                        </div>

                        <button type="submit" className="auth-button" disabled={loading}>
                            {loading ? <Loader2 className="spinner" /> : 'Reset Password'}
                        </button>

                        <div className="back-to-login">
                            <Link to="/login">
                                <ArrowLeft size={18} />
                                Kembali ke Login
                            </Link>
                        </div>
                    </form>
                )}
            </div>
        </div>
    );
};

export default ForgotPassword;
