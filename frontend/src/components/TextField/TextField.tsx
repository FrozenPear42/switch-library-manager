import {
  Input,
  Label,
  Text,
  TextField,
  TextFieldProps,
} from "react-aria-components";

import styles from "./TextField.module.css";

type AppTextFieldProps = TextFieldProps & {
  label?: string;
  description?: string;
};

export default function AppTextField({
  label,
  description,
  ...props
}: AppTextFieldProps) {
  return (
    <TextField {...props} className={styles.textField}>
      <Label className={styles.label}>{label}</Label>
      <Input className={styles.input} />
      {description && <Text slot="description">{description}</Text>}
    </TextField>
  );
}
