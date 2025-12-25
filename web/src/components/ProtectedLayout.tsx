import {Outlet} from 'react-router';

export default function ProtectedLayout() {

    return (
    <div className="flex h-screen flex-col bg-base-200">
        <Outlet />
    </div>
    );
}
