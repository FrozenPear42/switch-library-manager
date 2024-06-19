import {
  Modal,
  Dialog,
  ModalOverlay,
  DialogTrigger,
  Button,
} from "react-aria-components";

import styles from "./OrganizeFilesModal.module.css";
import Toggle from "../../components/Toggle/Toggle";
import AppTextField from "../../components/TextField/TextField";

type OrganizeFilesModalProps = {
  isOpened: boolean;
  onOpen: (isOpen: boolean) => void;
};

export default function OrganizeFilesModal({
  isOpened,
  onOpen,
}: OrganizeFilesModalProps) {
  console.log("xD", isOpened);

  return (
    <>
      <div>aaaaaaaa</div>
      <DialogTrigger>
        <Button>aaaa</Button>
        <ModalOverlay className={styles.overlay}>
          <Modal
            isDismissable={false}
            isKeyboardDismissDisabled={true}
            // onOpenChange={(isOpen) => onOpen(isOpen)}
            // isOpen={isOpened}
            className={styles.modal}
          >
            <Dialog className={styles.dialog}>
              <div className={styles.dialogHeader}>Organize files</div>
              <div className={styles.spacer}></div>
              <div>
                <div className={styles.formGroup}>
                  <label>Base directory</label>
                  <div>D:\Switch\Roms</div>
                </div>
                <div className={styles.formGroup}>
                  <Toggle label="Delete empty folders"></Toggle>
                </div>
                <div className={styles.formGroup}>
                  <Toggle label="Delete old updates"></Toggle>
                </div>
                <div className={styles.formGroup}>
                  <Toggle label="Create folders"></Toggle>
                </div>
                <div className={styles.formGroup}>
                  <AppTextField label="Folder name pattern"></AppTextField>
                </div>
                <div className={styles.formGroup}>
                  <Toggle label="Rename files"></Toggle>
                </div>
                <div className={styles.formGroup}>
                  <AppTextField label="File name pattern"></AppTextField>
                </div>
              </div>
              <button>Start</button>
            </Dialog>
          </Modal>
        </ModalOverlay>
      </DialogTrigger>
    </>
  );
}
