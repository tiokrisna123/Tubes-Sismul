import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Heart, Mail, Lock, User, Loader2 } from 'lucide-react';
import './Auth.css';

const Register = () => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { register } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (password !== confirmPassword) {
            setError('Password tidak cocok');
            return;
        }

        if (password.length < 6) {
            setError('Password minimal 6 karakter');
            return;
        }

        setLoading(true);

        try {
            await register(email, password, name);
            navigate('/onboarding');
        } catch (err) {
            setError(err.response?.data?.error || 'Registrasi gagal. Silakan coba lagi.');
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
                    <h1>Buat Akun Baru</h1>
                    <p>Mulai perjalanan kesehatan Anda bersama kami</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && <div className="auth-error">{error}</div>}

                    <div className="input-group">
                        <User className="input-icon" />
                        <input
                            type="text"
                            placeholder="Nama Lengkap"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            required
                        />
                    </div>

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
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <Lock className="input-icon" />
                        <input
                            type="password"
                            placeholder="Konfirmasi Password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                        />
                    </div>

                    <button type="submit" className="auth-button" disabled={loading}>
                        {loading ? <Loader2 className="spinner" /> : 'Daftar'}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>Sudah punya akun? <Link to="/login">Login di sini</Link></p>
                </div>
            </div>
        </div>
    );
};

export default Register;
