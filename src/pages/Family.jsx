import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { familyAPI } from '../services/api';
import {
    ArrowLeft, Users, UserPlus, Check, X, Heart,
    Loader2, Mail, Eye, Trash2, Bell
} from 'lucide-react';
import './Family.css';

const Family = () => {
    const [members, setMembers] = useState([]);
    const [requests, setRequests] = useState({ received: [], sent: [] });
    const [loading, setLoading] = useState(true);
    const [showInvite, setShowInvite] = useState(false);
    const [inviteEmail, setInviteEmail] = useState('');
    const [relationship, setRelationship] = useState('family');
    const [inviting, setInviting] = useState(false);
    const [error, setError] = useState('');
    const [activeTab, setActiveTab] = useState('members');

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const [membersRes, requestsRes] = await Promise.all([
                familyAPI.getMembers(),
                familyAPI.getRequests(),
            ]);
            setMembers(membersRes.data.data || []);
            setRequests(requestsRes.data.data || { received: [], sent: [] });
        } catch (err) {
            console.error('Failed to load family data:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleInvite = async (e) => {
        e.preventDefault();
        setError('');
        setInviting(true);

        try {
            await familyAPI.invite({ member_email: inviteEmail, relationship });
            setShowInvite(false);
            setInviteEmail('');
            loadData();
        } catch (err) {
            setError(err.response?.data?.error || 'Gagal mengirim undangan');
        } finally {
            setInviting(false);
        }
    };

    const handleApprove = async (id) => {
        try {
            await familyAPI.approve(id);
            loadData();
        } catch (err) {
            console.error('Failed to approve:', err);
        }
    };

    const handleReject = async (id) => {
        try {
            await familyAPI.reject(id);
            loadData();
        } catch (err) {
            console.error('Failed to reject:', err);
        }
    };

    const handleRemove = async (id) => {
        if (!confirm('Yakin ingin menghapus anggota keluarga ini?')) return;
        try {
            await familyAPI.remove(id);
            loadData();
        } catch (err) {
            console.error('Failed to remove:', err);
        }
    };

    const relationships = [
        { value: 'parent', label: 'Orang Tua' },
        { value: 'child', label: 'Anak' },
        { value: 'spouse', label: 'Pasangan' },
        { value: 'sibling', label: 'Saudara' },
        { value: 'family', label: 'Keluarga Lain' },
    ];

    if (loading) {
        return (
            <div className="family-loading">
                <Loader2 className="spinner" />
                <p>Memuat data keluarga...</p>
            </div>
        );
    }

    const pendingCount = requests.received?.length || 0;

    return (
        <div className="family-page">
            <header className="page-header">
                <Link to="/dashboard" className="back-btn">
                    <ArrowLeft size={20} />
                </Link>
                <h1>Keluarga</h1>
                <button className="invite-btn" onClick={() => setShowInvite(true)}>
                    <UserPlus size={18} />
                </button>
            </header>

            <div className="family-container">
                <div className="tabs">
                    <button
                        className={`tab ${activeTab === 'members' ? 'active' : ''}`}
                        onClick={() => setActiveTab('members')}
                    >
                        <Users size={18} /> Anggota
                    </button>
                    <button
                        className={`tab ${activeTab === 'requests' ? 'active' : ''}`}
                        onClick={() => setActiveTab('requests')}
                    >
                        <Bell size={18} /> Permintaan
                        {pendingCount > 0 && <span className="badge">{pendingCount}</span>}
                    </button>
                </div>

                {activeTab === 'members' && (
                    <div className="members-list">
                        {members.length === 0 ? (
                            <div className="empty-state">
                                <Users size={48} />
                                <h3>Belum Ada Anggota</h3>
                                <p>Undang keluarga untuk memantau kesehatan bersama</p>
                                <button onClick={() => setShowInvite(true)} className="cta-btn">
                                    <UserPlus size={18} /> Undang Sekarang
                                </button>
                            </div>
                        ) : (
                            members.map((member) => (
                                <div key={member.id} className="member-card">
                                    <div className="member-avatar">
                                        {member.member_name?.charAt(0).toUpperCase() || 'U'}
                                    </div>
                                    <div className="member-info">
                                        <h4>{member.member_name}</h4>
                                        <span className="member-relation">{member.relationship}</span>
                                    </div>
                                    <div className="member-actions">
                                        <Link to={`/family/${member.id}/health`} className="action-btn view">
                                            <Eye size={16} /> Lihat
                                        </Link>
                                        <button onClick={() => handleRemove(member.id)} className="action-btn remove">
                                            <Trash2 size={16} />
                                        </button>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                )}

                {activeTab === 'requests' && (
                    <div className="requests-section">
                        {requests.received?.length > 0 && (
                            <>
                                <h3>Undangan Masuk</h3>
                                {requests.received.map((req) => (
                                    <div key={req.id} className="request-card">
                                        <div className="request-info">
                                            <h4>{req.from_name}</h4>
                                            <span>{req.from_email}</span>
                                            <span className="relation-badge">{req.relationship}</span>
                                        </div>
                                        <div className="request-actions">
                                            <button onClick={() => handleApprove(req.id)} className="approve-btn">
                                                <Check size={18} />
                                            </button>
                                            <button onClick={() => handleReject(req.id)} className="reject-btn">
                                                <X size={18} />
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </>
                        )}

                        {requests.sent?.length > 0 && (
                            <>
                                <h3>Undangan Terkirim</h3>
                                {requests.sent.map((req) => (
                                    <div key={req.id} className="request-card sent">
                                        <div className="request-info">
                                            <h4>{req.member_email}</h4>
                                            <span className="status pending">Menunggu</span>
                                        </div>
                                    </div>
                                ))}
                            </>
                        )}

                        {(!requests.received?.length && !requests.sent?.length) && (
                            <div className="empty-state small">
                                <Bell size={32} />
                                <p>Tidak ada permintaan</p>
                            </div>
                        )}
                    </div>
                )}
            </div>

            {showInvite && (
                <div className="modal-overlay" onClick={() => setShowInvite(false)}>
                    <div className="modal" onClick={(e) => e.stopPropagation()}>
                        <h2>Undang Anggota Keluarga</h2>
                        <p>Kirim undangan ke email anggota keluarga yang sudah terdaftar</p>

                        <form onSubmit={handleInvite}>
                            {error && <div className="error-msg">{error}</div>}

                            <div className="form-group">
                                <label><Mail size={16} /> Email</label>
                                <input
                                    type="email"
                                    placeholder="email@example.com"
                                    value={inviteEmail}
                                    onChange={(e) => setInviteEmail(e.target.value)}
                                    required
                                />
                            </div>

                            <div className="form-group">
                                <label><Heart size={16} /> Hubungan</label>
                                <select
                                    value={relationship}
                                    onChange={(e) => setRelationship(e.target.value)}
                                >
                                    {relationships.map((rel) => (
                                        <option key={rel.value} value={rel.value}>{rel.label}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="modal-actions">
                                <button type="button" onClick={() => setShowInvite(false)} className="btn-cancel">
                                    Batal
                                </button>
                                <button type="submit" className="btn-submit" disabled={inviting}>
                                    {inviting ? <Loader2 className="spinner" /> : 'Kirim Undangan'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Family;
