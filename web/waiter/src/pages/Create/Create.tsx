import * as React from 'react';
import './Create.css';
import Header from '../../components/Header/Header';
import {RouteProps, useHistory} from 'react-router';
import {useForm} from 'react-hook-form';
import {useEffect, useState} from 'react';
import kitchenService from '../../services/kitchenService';
interface CreateTypes extends RouteProps {
    page: string
}

const Create: React.FC<CreateTypes> = (props: CreateTypes) => {
    const history = useHistory();
    const { register, handleSubmit, watch } = useForm();
    const onSubmit = (data: any) => console.log(data);
    const [firebaseProjects, setFirebaseProjects] = useState<any[]>([]);
    const back = () => {
        history.goBack();
    };

    const deployType = watch('deployType', 'local');

    useEffect(() => {
        kitchenService.requestFirebaseAccounts().then(res => {
            res.json().then(data => {
                const projectNames: any[] = [];
                data.results.forEach((result: any) => {
                    projectNames.push(result.name);
                });
                setFirebaseProjects(projectNames);
            });
        });
    }, []);

    return (
        <div className="container--create">
            <Header />
            <hr />
            <div className="container container--menu">
                <div className="menu">
                    <div className="menu__header">
                        <span>Cooking</span>
                        <p className="menu__header__subtitle">Create a new specification for your app and deploy it automatically.</p>
                    </div>
                    <div className="menu__content">
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <label>App Name*: </label><input {...register("name", { required: true})} />&nbsp;
                            <label>Color 1: </label><input {...register("color")} />&nbsp;
                            <label>Title 1: </label><input {...register("title")} />&nbsp;
                            <label>Deploy Type*: </label>
                            <select {...register("deployType", { required: true})}>
                                <option value="local">Local</option>
                                <option value="firebase">Firebase</option>
                            </select>
                            { deployType === 'firebase' &&
                                <div>
                                    <label>Select Project</label>
                                    <select>
                                        {firebaseProjects.map((project, index) => {
                                            return (
                                                <option key={index}>{project}</option>
                                            )
                                        })}
                                    </select>
                                    {/* If options < 1, ask the user to enter information to create a new firebase project */}
                                </div>
                            }
                            <input type="submit" />
                        </form>
                    </div>
                    <button className="button" onClick={back}>Back</button>
                </div>
            </div>
        </div>
    );
};

export default Create;
