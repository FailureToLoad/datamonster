import {Outlet} from 'react-router';
import styles from './DefaultLayout.module.css';
import {GlossaryProvider} from '~/context/glossary';

export default function Layout() {
    return (
        <GlossaryProvider>
            <div className={styles.layout}>
                <Outlet />
            </div>
        </GlossaryProvider>
    );
}
