import './Home.css';
import Header from "../../components/Header/Header";
import Menu from "../../components/Menu/Menu";

const Home = () => {
    return (
        <div className="container">
            <Header />
            <div className="container container--menu">
                <Menu />
            </div>
        </div>
    );
}

export default Home;