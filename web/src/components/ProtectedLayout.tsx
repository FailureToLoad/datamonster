import {Outlet} from 'react-router';

export default function ProtectedLayout() {

    return (
    <div className="flex h-screen flex-col items-center justify-center">
        <Outlet />
    </div>
    );
}
