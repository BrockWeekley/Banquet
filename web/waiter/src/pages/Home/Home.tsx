import './Home.css';
import Header from "../../components/Header/Header";
import Menu from "../../components/Menu/Menu";
import * as React from 'react';

const Home = () => {
    return (
        <div className="container">
            <Header />
            <div className="container container--menu">
                <Menu />
            </div>
        </div>
    );
};

export default Home;
