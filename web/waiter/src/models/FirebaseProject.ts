export interface FirebaseProject {
    displayName: string;
    name: string;
    projectId: string;
    projectNumber: string;
    resources: Resource;
    state: string;
}

interface Resource {
    hostingSite: string;
}
