import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Heart, Mail, Lock, Loader2, Eye, EyeOff } from 'lucide-react';
import './Auth.css';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            await login(email, password);
            navigate('/dashboard');
        } catch (err) {
            setError(err.response?.data?.error || 'Login failed. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-container">
            {/* Animated Background Elements */}
            <div className="auth-bg-elements">
                {/* Floating Particles */}
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>
                <div className="particle"></div>

                {/* Glowing Orbs */}
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>
                <div className="glow-orb"></div>

                {/* Twinkling Stars */}
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>
                <div className="star"></div>

                {/* Waves */}
                <div className="wave"></div>
                <div className="wave"></div>

                {/* Floating Circles */}
                <div className="floating-circle"></div>
                <div className="floating-circle"></div>
                <div className="floating-circle"></div>
                <div className="floating-circle"></div>
                <div className="floating-circle"></div>

                {/* Floating Health Bubbles */}
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>
                <div className="health-bubble"></div>

                {/* DNA Helix Animation */}
                <div className="dna-strand"></div>
                <div className="dna-strand"></div>
            </div>

            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <Heart className="logo-icon" />
                        <span>Live for Health</span>
                    </div>
                    <h1>Selamat Datang Kembali</h1>
                    <p>Login untuk melanjutkan ke dashboard kesehatan Anda</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && <div className="auth-error">{error}</div>}

                    <div className="input-group">
                        <Mail className="input-icon" />
                        <input
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <Lock className="input-icon" />
                        <input
                            type={showPassword ? "text" : "password"}
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
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

                    <div className="forgot-password">
                        <Link to="/forgot-password">Lupa Password?</Link>
                    </div>

                    <button type="submit" className="auth-button" disabled={loading}>
                        {loading ? <Loader2 className="spinner" /> : 'Login'}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>Belum punya akun? <Link to="/register">Daftar sekarang</Link></p>
                </div>
            </div>
        </div>
    );
};

export default Login;
