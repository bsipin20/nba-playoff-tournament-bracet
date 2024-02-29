import React, { useEffect, useState } from 'react';
import axios from '../axiosInstance'; // Import your axios instance with the interceptor

const Dashboard = () => {
    const [brackets, setBrackets] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {

        const fetchBrackets = async () => {
            try {
                const token = localStorage.getItem('authToken');

                if (!token) {
                    setError('No token available');
                    return;
                }

                const response = await axios.get('/brackets', {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });

                setBrackets(response.data);
            } catch (error) {
                console.error('Error fetching brackets:', error);
                setError('Error fetching brackets');
            }
        }
        fetchBrackets();
    }, []); // Empty dependency array to ensure the effect runs only once on component mount

    return (
        <div>
            {error ? (
                <p>Error: {error}</p>
            ) : (
                <ul>
                    {brackets.map((bracket, index) => (
                        <li key={index}>{bracket.name}</li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default Dashboard;