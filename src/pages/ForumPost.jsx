import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { forumAPI } from '../services/api';
import { ArrowLeft, Heart, MessageCircle, Send, Trash2, Loader2 } from 'lucide-react';
import './ForumPost.css';

const ForumPost = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState('');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        fetchPost();
    }, [id]);

    const fetchPost = async () => {
        try {
            const res = await forumAPI.getPost(id);
            setPost(res.data.post);
            setComments(res.data.comments || []);
        } catch (error) {
            console.error('Failed to fetch post:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleLike = async () => {
        try {
            const res = await forumAPI.toggleLike(id);
            setPost({ ...post, is_liked: res.data.is_liked, likes_count: res.data.likes_count });
        } catch (error) {
            console.error('Failed to toggle like:', error);
        }
    };

    const handleComment = async (e) => {
        e.preventDefault();
        if (!newComment.trim()) return;

        setSubmitting(true);
        try {
            const res = await forumAPI.addComment(id, { content: newComment });
            setComments([...comments, res.data]);
            setNewComment('');
            setPost({ ...post, comments_count: post.comments_count + 1 });
        } catch (error) {
            console.error('Failed to add comment:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    if (loading) {
        return (
            <div className="forum-post-loading">
                <Loader2 className="spinner" size={48} />
            </div>
        );
    }

    if (!post) {
        return (
            <div className="forum-post-loading">
                <p>Post tidak ditemukan</p>
                <button onClick={() => navigate('/forum')}>Kembali</button>
            </div>
        );
    }

    return (
        <div className="forum-post-page">
            <header className="forum-post-header">
                <button className="back-btn" onClick={() => navigate('/forum')}>
                    <ArrowLeft size={24} />
                </button>
                <h1>Diskusi</h1>
            </header>

            <div className="forum-post-container">
                <article className="post-detail">
                    <div className="post-header">
                        <div className="post-avatar">
                            {post.user_name?.charAt(0).toUpperCase() || 'U'}
                        </div>
                        <div className="post-meta">
                            <span className="post-author">{post.user_name}</span>
                            <span className="post-date">{formatDate(post.created_at)}</span>
                        </div>
                    </div>

                    <h2 className="post-title">{post.title}</h2>
                    <div className="post-content">
                        {post.content.split('\n').map((para, i) => (
                            <p key={i}>{para}</p>
                        ))}
                    </div>

                    <div className="post-stats">
                        <button
                            className={`stat-btn ${post.is_liked ? 'liked' : ''}`}
                            onClick={handleLike}
                        >
                            <Heart size={20} fill={post.is_liked ? '#ef4444' : 'none'} />
                            <span>{post.likes_count} suka</span>
                        </button>
                        <span className="stat-item">
                            <MessageCircle size={20} />
                            <span>{post.comments_count} komentar</span>
                        </span>
                    </div>
                </article>

                <section className="comments-section">
                    <h3>Komentar ({comments.length})</h3>

                    <form className="comment-form" onSubmit={handleComment}>
                        <input
                            type="text"
                            placeholder="Tulis komentar..."
                            value={newComment}
                            onChange={(e) => setNewComment(e.target.value)}
                        />
                        <button type="submit" disabled={submitting || !newComment.trim()}>
                            {submitting ? <Loader2 className="spinner" size={20} /> : <Send size={20} />}
                        </button>
                    </form>

                    <div className="comments-list">
                        {comments.length === 0 ? (
                            <p className="no-comments">Belum ada komentar. Jadilah yang pertama!</p>
                        ) : (
                            comments.map((comment) => (
                                <div key={comment.id} className="comment-item">
                                    <div className="comment-avatar">
                                        {comment.user_name?.charAt(0).toUpperCase() || 'U'}
                                    </div>
                                    <div className="comment-body">
                                        <div className="comment-header">
                                            <span className="comment-author">{comment.user_name}</span>
                                            <span className="comment-date">{formatDate(comment.created_at)}</span>
                                        </div>
                                        <p className="comment-text">{comment.content}</p>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </section>
            </div>
        </div>
    );
};

export default ForumPost;
