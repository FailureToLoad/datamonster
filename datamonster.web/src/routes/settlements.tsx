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
import api, { CreateSettlementRequest, Settlement } from "@/api/api";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "@/auth/auth-context";
import Spinner from "@/components/spinner";

const formSchema = z.object({
  settlementName: z.string().min(2).max(100),
});

interface CreateSettlementProps {
  token: string;
  settlements: Array<Settlement>;
}

function CreateSettlementModal(props: CreateSettlementProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      settlementName: "",
    },
  });
  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      const request: CreateSettlementRequest = {
        name: values.settlementName,
      };
      const response = await api.createSettlement(request, props.token);
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

function SettlementList({ settlements, token }: CreateSettlementProps) {
  if (settlements.length < 1) {
    return <CreateSettlementModal token={token} settlements={settlements} />;
  }

  const cards = settlements.map((settlement) => (
    <Card key={settlement.id}>
      <CardHeader>
        <CardTitle>{settlement.name}</CardTitle>
        <CardDescription>
          Settlement ID: {settlement.id} | Settlement Key: {settlement.name}
        </CardDescription>
      </CardHeader>
    </Card>
  ));
  //cards.push(<CreateSettlementModal token={token} settlements={settlements} />);
  return cards;
}

function Settlements() {
  const [settlements, setSettlements] = useState<Array<Settlement>>(
    Array<Settlement>(),
  );
  const [isLoading, setIsLoading] = useState(settlements.length < 1);
  const [token, setToken] = useState("");
  const { currentUser } = useContext(AuthContext);

  useEffect(() => {
    if (!currentUser) {
      return;
    }
    if (!token) {
      currentUser.getIdToken().then((idToken) => {
        setToken(idToken);
      });
    } else {
      api.getSettlementsForUser(token).then((val) => {
        setSettlements(val);
        setIsLoading(false);
      });
    }
  }, [currentUser, token]);

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      {isLoading ? (
        <Spinner />
      ) : (
        <SettlementList token={token} settlements={settlements} />
      )}
    </div>
  );
}

export default Settlements;
