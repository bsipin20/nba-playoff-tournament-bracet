import React, {useState , useEffect} from 'react';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { Link } from 'react-router-dom';

import { validate } from './validate';
import  {notify} from './toast';
import styles from './SignUp.module.css';

const Login = () => {
    const [data , setData] = useState({
        email:"",
        password:"",
    });
    const [errors , setErrors] = useState({});
    const [touched , setTouched] = useState({});
    const [token, setToken] = useState(null);


    useEffect(() => {
        setErrors(validate(data , "login"))
    }, [data, touched])

    // Event Handlers
    const changeHandler = event => {
            setData({...data , [event.target.name]:event.target.value})
    }

    const focusHandler = event => {
        setTouched({...touched , [event.target.name]: true})
    }

    const submitHandler = event => {
        event.preventDefault();
        if (!Object.keys(errors).length) {
            notify("You login successfully" , "success")
        } else {
            notify("Invalid data!" , "error")
            setTouched({
                email:true,
                password:true,
            })
        }
    }

    const postDataLogin = async () => {
        try {
            const response = await fetch(`http://localhost:8080/v1/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ "username": data.email, "password": data.password }),
            })
            console.log(response)
            if (response.ok) {
                const responseData = await response.json();
                const jwtToken = responseData.token;
                setToken(jwtToken);
                localStorage.setItem('authToken', jwtToken);
				localStorage.setItem('userId', responseData.userId);
                window.location.href = '/dashboard';
            } else {
                setErrors('Invalid credentials');
            }

        } catch (error) {
            console.error('Error posting data:', error);
        }
    };

    return (
        <div className={styles.container}>
            <form onSubmit={submitHandler} className={styles.formContainer}>
                <h2 className={styles.header}>Login</h2>
                <div className={styles.formField}>
                    <label>Email</label>
                    <input type="text" name="email" value={data.email} onChange={changeHandler} onFocus={focusHandler} className={(errors.email && touched.email) ? styles.uncompleted : styles.formInput} />
                    {errors.email && touched.email && <span>{errors.email}</span>}
                </div>
                <div className={styles.formField}>
                    <label>Password</label>
                    <input type="password" name="password" value={data.password} onChange={changeHandler} onFocus={focusHandler} className={(errors.password && touched.password) ? styles.uncompleted : styles.formInput} />
                    {errors.password && touched.password && <span>{errors.password}</span>}
                </div>
                <div className={styles.formButtons}>
                <Link to='/signup'>Signup</Link>                   
                    <button onClick={postDataLogin} type="submit">Login</button>
                </div>
            </form>
            <ToastContainer />
        </div>
    );
};

export default Login;
