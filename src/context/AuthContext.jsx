{/*import { createContext, useContext, useState, useEffect } from 'react';
import { authAPI } from '../services/api';

const AuthContext = createContext(null);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    return context;
};

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const savedUser = localStorage.getItem('user');

        if (token && savedUser) {
            setUser(JSON.parse(savedUser));
            // Verify token is still valid
            authAPI.getProfile()
                .then((res) => {
                    setUser(res.data.data);
                    localStorage.setItem('user', JSON.stringify(res.data.data));
                })
                .catch(() => {
                    localStorage.removeItem('token');
                    localStorage.removeItem('user');
                    setUser(null);
                })
                .finally(() => setLoading(false));
        } else {
            setLoading(false);
        }
    }, []);

    const login = async (email, password) => {
        const response = await authAPI.login({ email, password });
        const { token, user } = response.data.data;
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        setUser(user);
        return user;
    };

    const register = async (email, password, name) => {
        const response = await authAPI.register({ email, password, name });
        const { token, user } = response.data.data;
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        setUser(user);
        return user;
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setUser(null);
    };

    const updateUser = (userData) => {
        setUser(userData);
        localStorage.setItem('user', JSON.stringify(userData));
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, register, logout, updateUser }}>
            {children}
        </AuthContext.Provider>
    );
}; */}

import { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext(null);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    return context;
};

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // LOGIKA DEMO: Hilangkan pengecekan ke authAPI, langsung baca dari localStorage
    useEffect(() => {
        const savedUser = localStorage.getItem('user');
        if (savedUser) {
            setUser(JSON.parse(savedUser));
        }
        setLoading(false); // Langsung set loading selesai
    }, []);

    // LOGIKA DEMO: Bypass Login API
    const login = async (email, password) => {
        // Simulasi jeda loading 1 detik agar terasa seperti memanggil server (opsional untuk UX)
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Buat data pengguna palsu
        const mockUser = {
            id: 1,
            name: "Guest (Demo Mode)",
            email: email,
            role: "user"
        };

        localStorage.setItem('token', 'token-demo-rahasia-123');
        localStorage.setItem('user', JSON.stringify(mockUser));
        setUser(mockUser);
        return mockUser;
    };

    // LOGIKA DEMO: Bypass Register API
    const register = async (email, password, name) => {
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const mockUser = {
            id: 2,
            name: name || "New Guest",
            email: email,
            role: "user"
        };

        localStorage.setItem('token', 'token-demo-rahasia-123');
        localStorage.setItem('user', JSON.stringify(mockUser));
        setUser(mockUser);
        return mockUser;
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setUser(null);
    };

    const updateUser = (userData) => {
        setUser(userData);
        localStorage.setItem('user', JSON.stringify(userData));
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, register, logout, updateUser }}>
            {children}
        </AuthContext.Provider>
    );
};
