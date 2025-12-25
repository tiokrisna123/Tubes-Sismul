import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { forumAPI } from '../services/api';
import { ArrowLeft, Plus, Heart, MessageCircle, Loader2, Send, X } from 'lucide-react';
import './Forum.css';

const Forum = () => {
    const navigate = useNavigate();
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showNewPost, setShowNewPost] = useState(false);
    const [newPost, setNewPost] = useState({ title: '', content: '' });
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        fetchPosts();
    }, []);

    const fetchPosts = async () => {
        try {
            const res = await forumAPI.getPosts();
            setPosts(res.data || []);
        } catch (error) {
            console.error('Failed to fetch posts:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreatePost = async (e) => {
        e.preventDefault();
        if (!newPost.title.trim() || !newPost.content.trim()) return;

        setSubmitting(true);
        try {
            await forumAPI.createPost(newPost);
            setNewPost({ title: '', content: '' });
            setShowNewPost(false);
            fetchPosts();
        } catch (error) {
            console.error('Failed to create post:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleLike = async (postId) => {
        try {
            const res = await forumAPI.toggleLike(postId);
            setPosts(posts.map(p =>
                p.id === postId
                    ? { ...p, is_liked: res.data.is_liked, likes_count: res.data.likes_count }
                    : p
            ));
        } catch (error) {
            console.error('Failed to toggle like:', error);
        }
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'Baru saja';
        if (diffMins < 60) return `${diffMins} menit lalu`;
        if (diffHours < 24) return `${diffHours} jam lalu`;
        if (diffDays < 7) return `${diffDays} hari lalu`;
        return date.toLocaleDateString('id-ID');
    };

    if (loading) {
        return (
            <div className="forum-loading">
                <Loader2 className="spinner" size={48} />
                <p>Memuat forum...</p>
            </div>
        );
    }

    return (
        <div className="forum-page">
            <header className="forum-header">
                <button className="back-btn" onClick={() => navigate('/dashboard')}>
                    <ArrowLeft size={24} />
                </button>
                <h1>👥 Komunitas Sehat</h1>
                <button className="new-post-btn" onClick={() => setShowNewPost(true)}>
                    <Plus size={20} />
                </button>
            </header>

            <div className="forum-container">
                {posts.length === 0 ? (
                    <div className="empty-state">
                        <MessageCircle size={64} />
                        <h3>Belum Ada Diskusi</h3>
                        <p>Jadilah yang pertama memulai diskusi!</p>
                        <button className="cta-btn" onClick={() => setShowNewPost(true)}>
                            <Plus size={18} /> Buat Post
                        </button>
                    </div>
                ) : (
                    <div className="posts-list">
                        {posts.map((post) => (
                            <article
                                key={post.id}
                                className="post-card"
                                onClick={() => navigate(`/forum/${post.id}`)}
                            >
                                <div className="post-header">
                                    <div className="post-avatar">
                                        {post.user_name?.charAt(0).toUpperCase() || 'U'}
                                    </div>
                                    <div className="post-meta">
                                        <span className="post-author">{post.user_name}</span>
                                        <span className="post-date">{formatDate(post.created_at)}</span>
                                    </div>
                                </div>
                                <h3 className="post-title">{post.title}</h3>
                                <p className="post-preview">
                                    {post.content.length > 150
                                        ? post.content.substring(0, 150) + '...'
                                        : post.content}
                                </p>
                                <div className="post-actions">
                                    <button
                                        className={`action-btn like-btn ${post.is_liked ? 'liked' : ''}`}
                                        onClick={(e) => { e.stopPropagation(); handleLike(post.id); }}
                                    >
                                        <Heart size={18} fill={post.is_liked ? '#ef4444' : 'none'} />
                                        <span>{post.likes_count}</span>
                                    </button>
                                    <button className="action-btn">
                                        <MessageCircle size={18} />
                                        <span>{post.comments_count}</span>
                                    </button>
                                </div>
                            </article>
                        ))}
                    </div>
                )}
            </div>

            {/* New Post Modal */}
            {showNewPost && (
                <div className="modal-overlay" onClick={() => setShowNewPost(false)}>
                    <div className="modal" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h2>Buat Diskusi Baru</h2>
                            <button className="close-btn" onClick={() => setShowNewPost(false)}>
                                <X size={24} />
                            </button>
                        </div>
                        <form onSubmit={handleCreatePost}>
                            <div className="form-group">
                                <label>Judul</label>
                                <input
                                    type="text"
                                    placeholder="Apa yang ingin Anda diskusikan?"
                                    value={newPost.title}
                                    onChange={(e) => setNewPost({ ...newPost, title: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Konten</label>
                                <textarea
                                    placeholder="Tuliskan detail diskusi Anda..."
                                    value={newPost.content}
                                    onChange={(e) => setNewPost({ ...newPost, content: e.target.value })}
                                    rows={5}
                                    required
                                />
                            </div>
                            <button type="submit" className="submit-btn" disabled={submitting}>
                                {submitting ? <Loader2 className="spinner" size={20} /> : <Send size={20} />}
                                {submitting ? 'Memposting...' : 'Posting'}
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Forum;
