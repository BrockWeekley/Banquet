import * as React from 'react';
import './Create.css';
import Header from '../../components/Header/Header';
import {RouteProps, useHistory} from 'react-router';
import {useForm} from 'react-hook-form';
import {useEffect, useState} from 'react';
import kitchenService from '../../services/kitchenService';
import {FirebaseProject} from '../../models/FirebaseProject';
interface CreateTypes extends RouteProps {
    page: string
}

const Create: React.FC<CreateTypes> = (props: CreateTypes) => {
    const history = useHistory();
    const { register, handleSubmit, watch } = useForm();
    const onSubmit = (data: any) => {
        kitchenService.prepareCourse(data).then(res => {
            console.log(res);
        });
    };
    const [firebaseProjects, setFirebaseProjects] = useState<FirebaseProject[]>([]);
    const back = () => {
        history.goBack();
    };

    const deployType = watch('ProjectType', 'Local');

    useEffect(() => {
        kitchenService.requestFirebaseAccounts().then(res => {
            res.json().then(data => {
                const projects: FirebaseProject[] = [];
                data.results.forEach((result: any) => {
                    projects.push(result as FirebaseProject);
                });
                setFirebaseProjects(projects);
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
                    <div className="cooking__content">
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <label>App Name*: </label><input {...register("ProjectName", { required: true})} /><br /><br />
                            <label>GitHub URL*: </label><input {...register("GitHubURL", { required: true})} /><br /><br />
                            <label>Color 1: </label><input {...register("color")} /><br /><br />
                            <label>Title 1: </label><input {...register("title")} /><br /><br />
                            <label>Deploy Type*: </label>
                            <select {...register("ProjectType", { required: true})}>
                                <option value="Local">Local</option>
                                <option value="Web">Firebase Web</option>
                            </select><br /><br />
                            { deployType !== 'Local' &&
                                <div>
                                    <label>Select Project: </label>
                                    <select {...register("ParentID", { required: true})}>
                                        {firebaseProjects.map((project, index) => {
                                            return (
                                                <option key={index} value={project.projectNumber}>{project.name}</option>
                                            )
                                        })}
                                    </select>
                                    {/* If options < 1, ask the user to enter information to create a new firebase project */}
                                </div>
                            }
                            <br />
                            <input type="submit" />
                        </form>
                    </div><br /><br />
                    <button className="button" onClick={back}>Back</button>
                </div>
            </div>
        </div>
    );
};

export default Create;
