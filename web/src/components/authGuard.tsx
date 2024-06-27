import { withAuthenticationRequired } from "@auth0/auth0-react";
import React, { ComponentType } from "react";
import Spinner from "./spinner";

interface AuthenticationGuardProps {
  component: ComponentType;
}

const AuthGuard: React.FC<AuthenticationGuardProps> = ({ component }) => {
  const Component = withAuthenticationRequired(component, {
    onRedirecting: () => (
      <div className="page-layout">
        <Spinner />
      </div>
    ),
  });

  return <Component />;
};

export default AuthGuard;
