import "../.././App.css";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import Form from "./form";
import * as z from "zod";

import { RegisterDialogue } from "./register";
import { Authenticator } from "@/auth/authenticator";
import {
  ActionFunctionArgs,
  Navigate,
  redirect,
  useLoaderData,
  useNavigation,
} from "react-router-dom";
import { Label } from "@/components/ui/label";

const validator = z.object({
  username: z.string().min(3).max(50),
  password: z.string().min(6),
});

type FormData = z.infer<typeof validator>;

export async function LoginAction({ request }: ActionFunctionArgs) {
  const formData = Object.fromEntries(await request.formData()) as FormData;
  await Authenticator.signin(formData.username, formData.password);
  if (!Authenticator.isAuthenticated()) {
    return null;
  }

  return redirect(`/`);
}

function Login() {
  const { state } = useNavigation();
  const isAuthenticated = useLoaderData();
  if (isAuthenticated) {
    return <Navigate to="/select" />;
  }

  return (
    <div className="App">
      <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
        <div className="m-auto w-full bg-white lg:max-w-lg">
          <Card>
            <CardHeader className="space-y-1">
              <CardTitle className="text-center text-2xl">Sign in</CardTitle>
              <CardDescription className="text-center">
                Enter your username and password to login
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
              <Form validator={validator}>
                {(register, errors) => (
                  <>
                    <div className="space-y-2">
                      <div>
                        <Label htmlFor="username">Username</Label>
                        {errors.username && (
                          <p className="text-error text-sm italic">
                            {errors.username.message}
                          </p>
                        )}
                        <Input
                          id="username"
                          {...register("username")}
                          type="text"
                        />
                      </div>
                      <div>
                        <Label htmlFor="password">Password</Label>
                        {errors.password && (
                          <p className="text-error text-sm italic">
                            {errors.password.message}
                          </p>
                        )}
                        <Input
                          id="password"
                          {...register("password")}
                          type="password"
                        />
                      </div>

                      <div className="flex justify-between">
                        <Button type="submit" disabled={state === "submitting"}>
                          Sign In
                        </Button>
                        <RegisterDialogue />
                      </div>
                    </div>
                  </>
                )}
              </Form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

export default Login;
