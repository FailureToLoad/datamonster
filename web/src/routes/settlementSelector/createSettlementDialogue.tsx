import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  Form,
} from "@/components/ui/form";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Post } from "@/api/api";
import { Settlement } from "@/api/settlement";
import { useState } from "react";
import { Plus } from "lucide-react";
import { useAuth0 } from "@auth0/auth0-react";

const formSchema = z.object({
  settlementName: z.string().min(2).max(100),
});

type SettlementCreationRequest = {
  name: string;
};

export interface CreateSettlementProps {
  update: (s: Settlement) => void;
}
export function CreateSettlementDialogue() {
  const { getAccessTokenSilently } = useAuth0();
  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      settlementName: "",
    },
  });
  async function createSettlement(values: z.infer<typeof formSchema>) {
    try {
      const token = await getAccessTokenSilently();
      const request: SettlementCreationRequest = {
        name: values.settlementName,
      };
      await Post("settlement", request, token);
      setOpen(false);
    } catch (error) {
      console.log(error);
    }
  }
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          className="w-full"
          aria-label="launch add settlement dialogue"
        >
          <Plus className="h-6 w-6" />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form
            className="space-y-8"
            onSubmit={form.handleSubmit(createSettlement)}
          >
            <DialogHeader>
              <DialogTitle>Add Settlement</DialogTitle>
              <DialogDescription>Enter settlement details.</DialogDescription>
            </DialogHeader>

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
            <DialogFooter>
              <Button type="submit">Add</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
