import type {
  ListBoxItemProps,
  SelectProps,
  ValidationResult,
} from "react-aria-components";
import {
  FieldError,
  Text,
  ListBoxItem,
  Select,
  Label,
  Button,
  SelectValue,
  Popover,
  ListBox,
} from "react-aria-components";

import styles from "./Select.module.css";
import { IconChevronDown } from "@tabler/icons-react";

type AppSelectProps<T extends object> = Omit<SelectProps<T>, "children"> & {
  label?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
  items?: Iterable<T>;
  children: React.ReactNode | ((item: T) => React.ReactNode);
};

export function AppSelect<T extends object>({
  label,
  description,
  errorMessage,
  children,
  items,
  ...props
}: AppSelectProps<T>) {
  return (
    <Select
      {...props}
      className={styles.select}
      onSelectionChange={(k) => console.log(k)}
    >
      <Label className={styles.label}>{label}</Label>
      <Button className={styles.button}>
        <SelectValue className={styles.selectValue} />
        <IconChevronDown />
      </Button>
      {description && <Text slot="description">{description}</Text>}
      <FieldError>{errorMessage}</FieldError>
      <Popover className={styles.popover}>
        <ListBox items={items} className={styles.listBox}>
          {children}
        </ListBox>
      </Popover>
    </Select>
  );
}

export function AppSelectItem(props: ListBoxItemProps) {
  return <ListBoxItem {...props} className={styles.listItem} />;
}
