import { zodResolver } from "@hookform/resolvers/zod";
import { ReactNode } from "react";
import { FieldErrors, UseFormRegister, useForm } from "react-hook-form";
import { useSubmit } from "react-router-dom";
import { ZodType, z } from "zod";

interface Props<Validator extends ZodType<any, any, any>> {
  validator: Validator;
  children: (
    register: UseFormRegister<z.infer<Validator>>,
    errors: FieldErrors<z.infer<Validator>>,
  ) => ReactNode;
}

export default function Form<Validator extends ZodType<any, any, any>>(
  props: Props<Validator>,
) {
  const submit = useSubmit();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<Validator>>({ resolver: zodResolver(props.validator) });

  return (
    <form
      method="post"
      onSubmit={(event) => {
        const target = event.currentTarget;
        handleSubmit(() => {
          submit(target, { method: "post" });
        })(event);
      }}
    >
      {props.children(register, errors)}
    </form>
  );
}
