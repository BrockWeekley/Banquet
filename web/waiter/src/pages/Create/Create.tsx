import * as React from 'react';
import './Create.css';
import Header from '../../components/Header/Header';

const Create: React.FC = () => {
    return (
        <div className="container--create">
            <Header />
            <hr />
            Create a new dish
        </div>
    );
};

export default Create;
