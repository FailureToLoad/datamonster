import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  Form,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import api from "@/api/api";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "@/auth/auth-context";
import Spinner from "@/components/spinner";

const formSchema = z.object({
  settlementName: z.string().min(2).max(100),
  userId: z.string().min(6),
});

function CreateSettlementCard() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      settlementName: "",
      userId: "",
    },
  });
  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      const response = await api.createSettlement(values);
      console.log(response.name);
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>No Settlements Found</CardTitle>
        <CardDescription>
          Fill out the fields below to create one
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form className="space-y-8" onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="settlementName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Settlement Name</FormLabel>
                  <FormControl>
                    <Input type="text" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <Button type="submit">Create</Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}

function Settlements() {
  const [isLoading, setIsLoading] = useState(true);
  const { currentUser } = useContext(AuthContext);

  useEffect(() => {
    async function getSettlements() {
      if (!currentUser) {
        return;
      }
      const token = await currentUser?.getIdToken();
      const settlements = await api.getSettlementsForUser(token);
      console.log(settlements);
      setIsLoading(false);
    }
    getSettlements();
  }, []);

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      {isLoading ? <Spinner /> : <CreateSettlementCard />}
    </div>
  );
}

export default Settlements;
